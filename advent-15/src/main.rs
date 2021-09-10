use std::collections::HashMap;

fn main() {
    let starting_numbers = [0,6,1,7,2,19,20];
    let mut current_turn = 1;
    let mut last_said: HashMap<i32, i32> = HashMap::new();

    for starting_num in starting_numbers.iter() {
        last_said.insert(*starting_num, current_turn);
        current_turn += 1;
        println!("{}", starting_num);
    }

    let mut last_number = starting_numbers.last().unwrap().clone();
    let mut this_number = 0;

    last_said.remove(&last_number);

    loop {
        this_number = match last_said.get(&last_number) {
            Some(turn_last_said) => current_turn - turn_last_said - 1,
            None => 0,
        };

        println!("{}", 30000001 - current_turn);

        last_said.insert(last_number, current_turn - 1);


        current_turn += 1;
        last_number = this_number.clone();

        if current_turn == 30000001 {
            println!("{}", this_number);
            break;
        }
    }
}
