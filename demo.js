console.log("=====");
const startTime = Date.now();
let total = 0;
for (let i = 1; i <= 1000000; i++) {
    total += i;
}
const endTime = Date.now();
const elapsedMs = endTime - startTime;
console.log(`Time taken: ${elapsedMs.toFixed(3)} milliseconds`);
console.log(`Final result: ${total}`);

