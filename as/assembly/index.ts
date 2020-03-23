// The entry file of your WebAssembly module.
//import { wasi_abort } from "../node_modules/as-wasi/assembly/as-wasi";
import { Console } from "../node_modules/as-wasi/assembly";

export function add(a: i32, b: i32): i32 {
  Console.log("FOOOO");
  return a + b;
}
