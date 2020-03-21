// The entry file of your WebAssembly module.
import { wasi_abort } from "../node_modules/as-wasi/assembly/as-wasi";

export function add(a: i32, b: i32): i32 {
  return a + b;
}
