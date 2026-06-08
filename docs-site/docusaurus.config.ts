import { themes as prismThemes } from "prism-react-renderer";
import type { Config } from "@docusaurus/types";
import type * as Preset from "@docusaurus/preset-classic";

// This runs in Node.js - Don't use client-side code here (browser APIs, JSX...)

// 🔧 Concept-test CTA target. Swap this for your real form (Tally / Google Form
// / Typeform). Used by the navbar button, the landing hero and the footer.
const REQUEST_ACCESS_URL =
  "https://docs.google.com/forms/d/e/1FAIpQLSdB6WkNa1itYFSBwoAbvzFLrgn5ExsiblC-NviUhET78kU0zQ/viewform";
const GITHUB_URL =
  "https://github.com/Sreenivas-Sadhu-Prabhakara/business-launch-orchestrator";

const config: Config = {
  title: "Business Launch Orchestrator",
  tagline:
    "Launch a company in India, the Philippines or the US — one orchestrated flow.",
  favicon: "img/favicon.ico",

  future: {
    v4: true,
  },

  // GitHub Pages
  url: "https://sreenivas-sadhu-prabhakara.github.io",
  baseUrl: "/business-launch-orchestrator/",
  organizationName: "Sreenivas-Sadhu-Prabhakara",
  projectName: "business-launch-orchestrator",
  trailingSlash: false,

  onBrokenLinks: "warn",

  i18n: {
    defaultLocale: "en",
    locales: ["en"],
  },

  headTags: [
    {
      tagName: "link",
      attributes: { rel: "preconnect", href: "https://fonts.googleapis.com" },
    },
    {
      tagName: "link",
      attributes: {
        rel: "preconnect",
        href: "https://fonts.gstatic.com",
        crossorigin: "anonymous",
      },
    },
  ],

  stylesheets: [
    "https://fonts.googleapis.com/css2?family=Fraunces:opsz,wght@9..144,300..600&family=Hanken+Grotesk:wght@400..700&display=swap",
  ],

  customFields: {
    requestAccessUrl: REQUEST_ACCESS_URL,
    githubUrl: GITHUB_URL,
  },

  presets: [
    [
      "classic",
      {
        docs: {
          sidebarPath: "./sidebars.ts",
          routeBasePath: "docs",
          editUrl: `${GITHUB_URL}/tree/main/docs-site/`,
        },
        blog: false,
        theme: {
          customCss: "./src/css/custom.css",
        },
      } satisfies Preset.Options,
    ],
  ],

  themeConfig: {
    image: "img/docusaurus-social-card.jpg",
    colorMode: {
      defaultMode: "light",
      respectPrefersColorScheme: true,
    },
    navbar: {
      title: "Launch Orchestrator",
      items: [
        { to: "/docs/overview", label: "Docs", position: "left" },
        { to: "/docs/the-flow", label: "The flow", position: "left" },
        { to: "/docs/deploy-serverless", label: "Deploy", position: "left" },
        { href: GITHUB_URL, label: "GitHub", position: "right" },
        {
          href: REQUEST_ACCESS_URL,
          label: "Request access",
          position: "right",
          className: "navbar-cta",
        },
      ],
    },
    footer: {
      style: "dark",
      links: [
        {
          title: "Product",
          items: [
            { label: "Overview", to: "/docs/overview" },
            { label: "The 11-step flow", to: "/docs/the-flow" },
            { label: "Country coverage", to: "/docs/country-coverage" },
          ],
        },
        {
          title: "Build",
          items: [
            { label: "Architecture", to: "/docs/architecture" },
            { label: "API reference", to: "/docs/api-reference" },
            { label: "Deploy serverless", to: "/docs/deploy-serverless" },
          ],
        },
        {
          title: "More",
          items: [
            { label: "GitHub", href: GITHUB_URL },
            { label: "Request access", href: REQUEST_ACCESS_URL },
            { label: "FAQ", to: "/docs/faq" },
          ],
        },
      ],
      copyright: `© ${new Date().getFullYear()} Business Launch Orchestrator · Reference implementation — see the disclaimer. Built with Docusaurus.`,
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
      additionalLanguages: ["bash", "go", "json", "yaml", "sql"],
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
