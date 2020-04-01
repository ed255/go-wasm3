#!/usr/bin/env node

const fs = require("fs");
const {stringifyBigInts, unstringifyBigInts} = require("snarkjs");
const WitnessCalculatorBuilder = require("./witness_calculator.js");

const wasmName = "mycircuit2.wasm"
const inputName = "mycircuit2-input.json"

async function run () {
  const wasm = await fs.promises.readFile(wasmName);
  const input = unstringifyBigInts(JSON.parse(await fs.promises.readFile(inputName, "utf8")));

  console.log("input:", input);
  let options;
  const wc = await WitnessCalculatorBuilder(wasm, options);

  const w = await wc.calculateWitness(input);

  console.log("witness:", stringifyBigInts(w))

  // await fs.promises.writeFile(witnessName, JSON.stringify(stringifyBigInts(w), null, 1));

}

run().then(() => {
    process.exit();
});
