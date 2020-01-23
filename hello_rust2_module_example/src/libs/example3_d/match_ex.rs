pub fn match_fn1() {
    let v = vec!["one", "two", "three", "four", "five"];

    for (num, word) in v.iter().enumerate() {
        match num {
            0 => println!("_This is {}", num),
            1 => println!("_One of the solutions is {}", num),
            _ => println!("_Don't care about the rest "),
        }

        match word {
            &"five" => println!("* Yay this is five {}", word),
            _ => {}
        }
    }

    // or you could use 1..101
    for nr in 1..=100 {
        match nr {
            1 => println!("-> This is one {}", nr),
            2 => println!("-> This is two {}", nr),
            90..95 => println!("-> Between 90 and 95 {}", nr),
            _ => {}
        }
    }
}
