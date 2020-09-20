module.exports = {
  async rewrites() {
    return [
      {
        source: "/api/:path",
        destination: "http://localhost:8080/api/:path", // Matched parameters can be used in the destination
      },
    ];
  },
};
