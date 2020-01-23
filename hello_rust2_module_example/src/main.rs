#![allow(dead_code)]
#![allow(unused_imports)]
#![feature(exclusive_range_pattern)]

mod libs;
// or simply use libs::...
use crate::libs::example1_d::ex1_f as ex1;
use crate::libs::example2_d::person_d::person_1_f as ex2_p1;
use crate::libs::example2_d::person_d::person_2_f::{fn_person2, Person2};
use crate::libs::example3_d::match_ex as mex;

// use primal;
// use rand;

fn main() {
    // ex1::concurency1();
    // mex::match_fn1();
    println!("End program");
}
