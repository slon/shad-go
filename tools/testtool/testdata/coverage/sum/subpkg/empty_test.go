// This is subpkg package comment.
package subpkg

// Incorrect coverage comments:

// min coverage: . -1%

// min coverage: . 100.001%

// min coverage: . 100 %

// min coverage:. 10%

//  min coverage: . 19%

// min coverage: 90%

// Correct coverage comment:

// min coverage: . 90%

// Testtool uses first matching comment.

// min coverage: . 91%
