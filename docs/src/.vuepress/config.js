const { description } = require("../../package");

module.exports = {
  port: 3000,
  title: "Hetty",
  description: description,
  head: [
    ["meta", { name: "theme-color", content: "#30e3b7" }],
    ["meta", { name: "apple-mobile-web-app-capable", content: "yes" }],
    [
      "meta",
      { name: "apple-mobile-web-app-status-bar-style", content: "black" },
    ],
    [
      "meta",
      {
        property: "og:title",
        content: "Hetty",
      },
    ],
    [
      "meta",
      {
        property: "og:description",
        content: "An HTTP toolkit for security research.",
      },
    ],
    [
      "meta",
      {
        property: "og:image",
        content: "https://hetty.xyz/assets/hetty_v0.2.0_header.png",
      },
    ],
  ],
  themeConfig: {
    repo: "dstotijn/hetty",
    editLinks: true,
    docsDir: "docs/src",
    editLinkText: "",
    lastUpdated: true,
    logo: "/assets/logo.png",
    nav: [
      {
        text: "Guide",
        link: "/guide/",
      },
      {
        text: "Appendix",
        link: "/appendix/",
      },
    ],
    sidebar: {
      "/guide/": [
        {
          title: "Guide",
          collapsable: false,
          children: ["", "getting-started", "modules"],
        },
      ],
      "/appendix/": [
        {
          title: "Appendix",
          collapsable: false,
          children: [""],
        },
      ],
    },
  },
  plugins: ["@vuepress/plugin-back-to-top", "@vuepress/plugin-medium-zoom"],
  markdown: {
    toc: { includeLevel: [2] },
    extendMarkdown: (md) => {
      md.use(require("markdown-it-imsize"));
    },
  },
};
