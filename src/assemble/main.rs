use clap::Parser;
use scotus::db;
use std::error::Error;
use std::io;
use std::path::Path;

#[derive(Parser, Debug)]
struct Flags {
	#[clap(long, short, default_value_t = String::from("data"))]
	data_dir: String,

	#[clap(long, short, default_value_t = String::from(scotus::db::DEFAULT_CASES_URL))]
	scotusdb_cases_url: String,

	#[clap(long, short)]
	reset_data: bool,
}

fn ensure_data_dir<P: AsRef<Path>>(path: P, reset: bool) -> io::Result<()> {
	let path = path.as_ref();
	if path.exists() {
		if !reset {
			return Ok(());
		}
		std::fs::remove_dir_all(path)?;
	}
	std::fs::create_dir_all(path)
}

fn main() -> Result<(), Box<dyn Error>> {
	let flags = Flags::parse();

	ensure_data_dir(&flags.data_dir, flags.reset_data)?;

	let client =
		db::Client::with_data_dir(&flags.data_dir).with_cases_url(&flags.scotusdb_cases_url);

	let terms = client.read_terms()?;
	println!("{:?}", terms);

	Ok(())
}
