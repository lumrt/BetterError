// test.js
let x = 10;
let y = 0;
let result = x / y;

console.log("Result:", result);

console.log("Value of z:", z);

// typeError.js
const obj = undefined;
console.log(obj.property);

// referenceError.js
console.log(nonExistentVariable);

// syntaxError.js
const a = 5
const b = 10
console.log(a + b);

// unexpectedToken.js
const obj = { a: 1, b: 2, c: 3 };
const jsonString = JSON.stringify(obj;
console.log(jsonString);

// notAFunction.js
const obj = { a: 1, b: 2, c: 3 };
obj(); // TypeError: obj is not a function

// undefinedProperties.js
let undefinedVar;
console.log(undefinedVar.length); 
