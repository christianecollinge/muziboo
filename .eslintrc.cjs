module.exports = {
	root: true,
	env: {
		node: true,
		browser: true,
		es2022: true,
	},
	// Build output and generated files are not source; linting them produces
	// thousands of meaningless errors from minified bundles.
	ignorePatterns: [
		"node_modules/",
		"dist/",
		".astro/",
		".output/",
		".wrangler/",
		"public/",
	],
	extends: [
		"eslint:recommended",
		"plugin:@typescript-eslint/recommended",
		"plugin:astro/recommended",
	],
	parser: "@typescript-eslint/parser",
	parserOptions: {
		ecmaVersion: "latest",
		sourceType: "module",
	},
	plugins: ["@typescript-eslint"],
	rules: {
		"@typescript-eslint/no-unused-vars": [
			"warn",
			{
				argsIgnorePattern: "^_",
				varsIgnorePattern: "^_",
			},
		],
		"@typescript-eslint/no-explicit-any": "warn",
	},
	overrides: [
		// React rules apply only to React components. Applying them repo-wide
		// makes them fire on Astro templates, where `class` is correct.
		{
			files: ["*.jsx", "*.tsx"],
			extends: ["plugin:react/recommended", "plugin:react-hooks/recommended"],
			plugins: ["react"],
			parserOptions: {
				ecmaFeatures: {
					jsx: true,
				},
			},
			rules: {
				"react/react-in-jsx-scope": "off",
				"react/prop-types": "off",
			},
		},
		{
			files: ["*.astro"],
			parser: "astro-eslint-parser",
			parserOptions: {
				parser: "@typescript-eslint/parser",
				extraFileExtensions: [".astro"],
			},
		},
		// eslint-plugin-astro extracts inline <script> blocks into virtual files,
		// so rules for them must target this pattern rather than "*.astro".
		{
			files: ["**/*.astro/*.js", "*.astro/*.js"],
			parser: "espree",
			parserOptions: {
				ecmaVersion: "latest",
				sourceType: "module",
			},
			env: {
				browser: true,
				node: false,
			},
			// Globals installed at runtime by the analytics snippets.
			globals: {
				posthog: "readonly",
				dataLayer: "writable",
				gtag: "readonly",
			},
			rules: {
				// Block-scoped function declarations are valid ES2015+; the rule
				// only matters for ES5 and trips on the verbatim gtag snippet.
				"no-inner-declarations": "off",
			},
		},
		// Node scripts run outside the browser and may use CommonJS.
		{
			files: ["*.cjs", "scripts/**/*.js"],
			env: {
				node: true,
				browser: false,
			},
			rules: {
				"@typescript-eslint/no-var-requires": "off",
			},
		},
	],
	settings: {
		react: {
			version: "detect",
		},
	},
};
