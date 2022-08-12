const path = require('path');

module.exports = {
	entry: "./src/index.ts",
	entry: {
		"a/index": "./src/a/index.ts",
		"b/index": "./src/b/index.ts",
	},
	output: {
		filename: "[name].js",
		path: path.resolve(__dirname, "pkg/web/dist"),
	},
	...(process.env.production || !process.env.development
		? {}
		: { devtool: 'inline-source-map' }),
	resolve: {
		extensions: ['.ts', '.tsx', '.js'],
	},
	module: {
		rules: [
			{
				test: /\.tsx?$/,
				use: "ts-loader"
			},
			{
				test: /\.s?css$/,
				use: [
					"style-loader",
					"css-loader",
					"sass-loader",
				]
			}
		]
	},
	plugins: [],
	devServer: {
		static: './dist'
	},
};
