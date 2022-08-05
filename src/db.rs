use linked_hash_map::{Entry, LinkedHashMap};
use std::collections::HashMap;
use std::error::Error;
use std::fs::File;
use std::io;
use std::path::{Path, PathBuf};
use zip::ZipArchive;

pub const DEFAULT_CASES_URL: &'static str =
	"http://scdb.wustl.edu/_brickFiles/2021_01/SCDB_2021_01_justiceCentered_Citation.csv.zip";

const CASES_FILE_NAME: &'static str = "SCDB_justiceCentered_Citation.csv.zip";

struct Fields {
	term: usize,
	case_id: usize,
}

impl Fields {
	fn from_reader<R>(r: &mut csv::Reader<R>) -> Result<Fields, Box<dyn Error>>
	where
		R: io::Read,
	{
		let headers = r
			.headers()?
			.iter()
			.enumerate()
			.map(|(i, c)| (c, i))
			.collect::<HashMap<&str, usize>>();
		let term = match headers.get("term") {
			Some(i) => *i,
			None => Err("type column missing")?,
		};
		Ok(Fields { term })
	}
}

#[derive(Debug)]
pub struct Term {
	year: u32,
	cases: Vec<Case>,
}

impl Term {
	fn from_year(year: u32) -> Term {
		Term {
			year,
			cases: Vec::new(),
		}
	}

	fn from_byte_record<'a>(
		fields: &Fields,
		terms: &'a mut LinkedHashMap<u32, Term>,
		record: &csv::ByteRecord,
	) -> Result<&'a mut Term, Box<dyn Error>> {
		let year = String::from_utf8_lossy(record.get(fields.term).unwrap()).parse::<u32>()?;
		match terms.entry(year) {
			Entry::Occupied(e) => Ok(e.into_mut()),
			Entry::Vacant(e) => Ok(e.insert(Term::from_year(year))),
		}
	}
}

#[derive(Debug)]
pub struct Case {
	id: String,
	name: String,
	majority_votes: u8,
	minority_votes: u8,
	// add votes
	// add date
}

pub struct Client {
	data_dir: PathBuf,
	cases_url: String,
}

impl Client {
	pub fn with_data_dir<P: AsRef<Path>>(path: P) -> Client {
		Client {
			data_dir: PathBuf::from(path.as_ref()),
			cases_url: String::from(DEFAULT_CASES_URL),
		}
	}

	pub fn with_cases_url(mut self, url: &str) -> Client {
		self.cases_url = String::from(url);
		return self;
	}

	pub fn read_terms(&self) -> Result<Vec<Term>, Box<dyn Error>> {
		let dst = self.data_dir.join(CASES_FILE_NAME);
		if !dst.exists() {
			reqwest::blocking::get(&self.cases_url)?.copy_to(&mut File::create(&dst)?)?;
		}
		read_terms(File::open(&dst)?)
	}
}

fn read_terms(r: impl io::Read + io::Seek) -> Result<Vec<Term>, Box<dyn Error>> {
	let mut z = ZipArchive::new(r)?;
	let mut file = z.by_index(0)?;
	let mut c = csv::Reader::from_reader(&mut file);

	let fields = Fields::from_reader(&mut c)?;

	let mut terms = LinkedHashMap::new();

	for record in c.byte_records() {
		let _term = Term::from_byte_record(&fields, &mut terms, &record?)?;
	}

	Ok(terms.into_iter().map(|(_, v)| v).collect())
}
