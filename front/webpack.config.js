const path = require("path");

module.exports = {
  entry: "./src/index.ts", // Entry point of your application
  output: {
    path: path.resolve(__dirname, "dist"), // Output directory
    filename: "bundle.js", // Output bundle filename
  },
  resolve: {
    extensions: [".ts", ".tsx", ".js"], // File extensions to resolve
  },
  module: {
    rules: [
      {
        test: /\.ts?$/,
        use: "ts-loader",
        exclude: /node_modules/,
      },
      {
        test: /\.css$/i,
        use: ["style-loader", "css-loader"],
      },
    ],
  },
  devServer: {
    contentBase: path.join(__dirname, "public"), // Serve static files from the public directory
    port: 3000, // Port for the development server
  },
};
