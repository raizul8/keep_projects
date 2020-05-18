#![allow(dead_code)]
#![allow(unused_imports)]
#![feature(exclusive_range_pattern)]

mod mlib {
    include!("../mlibs/mlib_file1.rs");
}

mod libs;
use crate::libs::example1_d::ex1_f as ex1;
use crate::libs::example2_d::person_d::person_1_f as ex2_p1;
// or simply use libs::...
use libs::example2_d::person_d::person_2_f::{fn_person2, Person2};
use libs::example3_d::match_ex as mex;

// use primal;
// use rand;

fn main() {
    // ex1::concurency1();
    // mex::match_fn1(2);
    // println!("mlib: {}", mlib::fc1_from_mlib_file1(3));
    println!("End program");
}
