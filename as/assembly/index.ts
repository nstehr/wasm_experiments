// The entry file of your WebAssembly module.
//import { wasi_abort } from "../node_modules/as-wasi/assembly/as-wasi";
import { Console } from "../node_modules/as-wasi/assembly";


declare function stringFromHost(): usize;

export function add(a: i32, b: i32): i32 {
  return a + b;
}

// Allocate memory for a new byte array of
// size `len` and return the offset into
// the module's linear memory to the start
// of the block.
export function alloc(len: i32): usize {
  // create a new AssemblyScript byte array
  let buf = new Array<u8>(len);
  let buf_ptr = memory.data(8);
  // create a pointer to the byte array and
  // return it
  store<Array<u8>>(buf_ptr, buf);
  return buf_ptr;
}

export function array_sum(buf_ptr: usize, len: i32): u8 {
  let result: u8 = 0;
  for(let i = 0; i < len; i++) {
    result += load<u8>(buf_ptr + i) as u8;
  }
  return result as u8;
}

export function returnString(): ArrayBuffer {
  
  // sticking this here to test calling a host provided function
  // that modifies memory
  let strBuf = new Array<i32>(40);
  let buf_ptr = stringFromHost();
  for(let i = 0; i < 40; i++) {
    strBuf[i] = load<u8>(buf_ptr + i) as i32;
  }
  let val = String.fromCharCodes(strBuf);
  Console.log(val);

  let test = "This is a test of the emergency broadcast system";
  let str = String.UTF8.encode(test, true)
  return str;
}

