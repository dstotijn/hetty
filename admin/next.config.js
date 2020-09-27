const withCSS = require("@zeit/next-css");
const MonacoWebpackPlugin = require("monaco-editor-webpack-plugin");

module.exports = withCSS({
  trailingSlash: true,
  async rewrites() {
    return [
      {
        source: "/api/:path/",
        destination: "http://localhost:8080/api/:path/",
      },
    ];
  },
  webpack: (config) => {
    config.module.rules.push({
      test: /\.(png|jpg|gif|svg|eot|ttf|woff|woff2)$/,
      use: {
        loader: "url-loader",
        options: {
          limit: 100000,
        },
      },
    });

    config.plugins.push(
      new MonacoWebpackPlugin({
        languages: ["html", "json", "javascript"],
        filename: "static/[name].worker.js",
      })
    );

    return config;
  },
});
