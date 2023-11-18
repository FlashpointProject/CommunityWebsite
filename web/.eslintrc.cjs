/* eslint-env node */
module.exports = {
  extends: ['eslint:recommended', 'plugin:@typescript-eslint/recommended', 'plugin:react/jsx-runtime'],
  parser: '@typescript-eslint/parser',
  plugins: ['@typescript-eslint', 'react'],
  root: true,
  rules: {
    'eol-last': ['error', 'always'], // Enforce newline at the end of files
    'quotes': ['error', 'single'], // Enforce the use of single quotes
    '@typescript-eslint/quotes': ['error', 'single'], // Enforce the use of single quotes in TypeScript files
    '@typescript-eslint/semi': ['error', 'always'], // Enforce the use of semicolons in TypeScript files
    '@typescript-eslint/indent': ['error', 2, {
      'MemberExpression': 0,
      'SwitchCase': 1
    }],
    '@typescript-eslint/no-empty-function': ['error', { 'allow': ['arrowFunctions'] }],
    '@typescript-eslint/no-unused-vars': ['error', {
      'vars': 'all',
      'args': 'none'
    }],
    '@typescript-eslint/no-explicit-any': 'off',
    'react/no-direct-mutation-state': 'error',
    'react/no-render-return-value': 'error',
    'no-extra-semi': 'error',
    'no-trailing-spaces': 'error',
    'no-whitespace-before-property': 'error',
    'spaced-comment': 'error',
    'no-extra-boolean-cast': 'off'
  },
};
