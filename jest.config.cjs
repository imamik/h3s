module.exports = {
  testEnvironment: 'node',
  roots: ['<rootDir>/test/unit', '<rootDir>/test/integration', '<rootDir>/test/e2e'],
  setupFilesAfterEnv: ['<rootDir>/test/setup.js'],
  testPathIgnorePatterns: ['/node_modules/'],
  collectCoverage: true,
  coverageDirectory: 'coverage',
};
