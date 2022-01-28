// @ts-check

/**
 * @type {import('next').NextConfig}
 **/
const nextConfig = {
  reactStrictMode: true,
  trailingSlash: true,
  async rewrites() {
    return [
      {
        source: "/api/:path/",
        destination: "http://localhost:8080/api/:path/",
      },
    ];
  },
};

module.exports = nextConfig;
