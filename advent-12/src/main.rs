use std::fs;

fn part1() {
    // read a whole input string from the file
    let contents = fs::read_to_string("puzzle-input.txt").expect("Something went wrong reading the file");
    // split this up into individual instructions
    let instructions_read: Vec<&str> = contents.split("\n").collect();
    // position of ship (x, y)
    let mut pos: [i32; 2] = [0; 2];
    // direction ship is facing (x, y)
    let mut direction: [i32; 2] = [1, 0];

    // process each instruction one by one
    for instruction in instructions_read.iter() {
        let mut chars = instruction.chars();
        
        // this tuple stores the command in processed format
        let mut command = (chars.next().unwrap(), 0_i32);
        // temp String
        let mut temp_read = String::new();
        // read each char one by one until there are no chars left
        loop {
            let read = chars.next();
            if read == None {
                break;
            } else {
                temp_read += &String::from(read.unwrap());
            }
        }
        // convert the temp String to a number for the command tuple
        command.1 = temp_read.parse().unwrap();

        // carry out command
        match command.0 {
            // change the position by required amount
            'N' => pos[1] += command.1,
            'S' => pos[1] -= command.1,
            'E' => pos[0] += command.1,
            'W' => pos[0] -= command.1,
            // go in our stored direction a certain number of times (for both coordinates)
            'F' => {
                pos[0] += command.1 * direction[0];
                pos[1] += command.1 * direction[1];
            }
            // turn left however many times 90 fits into command.1
            'L' => {
                for _i in 0..(command.1 / 90) {
                    direction = [direction[1] * -1, direction[0]];
                }
            }
            // turn right however many times 90 fits into command.1
            'R' => {
                for _i in 0..(command.1 / 90) {
                    direction = [direction[1], direction[0] * -1];
                }
            }
            // unexpected command
            _ => panic!()
        }
    }
    // get manhattan distance
    println!("{}", pos[0].abs() + pos[1].abs());
}

fn part2() {
    // read a whole input string from the file
    let contents = fs::read_to_string("puzzle-input.txt").expect("Something went wrong reading the file");
    // split this up into individual instructions
    let instructions_read: Vec<&str> = contents.split("\n").collect();
    // position of ship (x, y)
    let mut pos: [i32; 2] = [0; 2];
    // position of waypoint ahead of the ship (x, y)
    let mut waypoint: [i32; 2] = [10, 1];

    // process each instruction one by one
    for instruction in instructions_read.iter() {
        let mut chars = instruction.chars();
        
        // this tuple stores the command in processed format
        let mut command = (chars.next().unwrap(), 0_i32);
        // temp String
        let mut temp_read = String::new();
        // read each char one by one until there are no chars left
        loop {
            let read = chars.next();
            if read == None {
                break;
            } else {
                temp_read += &String::from(read.unwrap());
            }
        }
        // convert the temp String to a number for the command tuple
        command.1 = temp_read.parse().unwrap();

        // carry out command
        match command.0 {
            // change position of waypoint by required amount
            'N' => waypoint[1] += command.1,
            'S' => waypoint[1] -= command.1,
            'E' => waypoint[0] += command.1,
            'W' => waypoint[0] -= command.1,
            // go forward by the required amount (multiple of distance of waypoint)
            'F' => {
                pos[0] += command.1 * waypoint[0];
                pos[1] += command.1 * waypoint[1];
            }
            // turn waypoint left however many times 90 fits into command.1
            'L' => {
                for _i in 0..(command.1 / 90) {
                    waypoint = [waypoint[1] * -1, waypoint[0]];
                }
            }
            // turn waypoint right however many times 90 fits into command.1
            'R' => {
                for _i in 0..(command.1 / 90) {
                    waypoint = [waypoint[1], waypoint[0] * -1];
                }
            }
            _ => panic!()
        }
    }
    // get manhattan distance
    println!("{}", pos[0].abs() + pos[1].abs());
}

fn main() {
    part1();
    part2();
}