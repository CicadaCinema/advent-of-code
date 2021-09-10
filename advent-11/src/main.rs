use std::fs;

// this function takes row and column indexes as input
// it returns the adjacent indexes as output, as long as they are within the correct range (see dimensions in main program)
fn get_adjacent(row: i32, col: i32) -> Vec<[i32; 2]> {
    let mut out_vector: Vec<[i32; 2]> = vec![];

    // first, naively add all the cells where each coordinate is at most 1 from the input
    for row_index in (row-1)..(row+2) {
        for col_index in (col-1)..(col+2) {
            out_vector.push([row_index, col_index]);
        }
    }

    // remove the starting coords (they appear in the middle) - this is precisely the input coords
    out_vector.remove(4);

    // go through the coords in reverse order (this way we can remove elements without causing havoc)
    for vec_index in (0..8).rev() {
        // if these coords do not fall within the desired range, remove them from the output
        if out_vector[vec_index][0]<0 || out_vector[vec_index][0]>93 || out_vector[vec_index][1]<0 || out_vector[vec_index][1]>91 {
            out_vector.remove(vec_index);
        }
    }

    out_vector
}

/*
Key for seat_array:
-1 is a floor cell
0 is an empty seat
1 is a filled seat
*/

fn part1() {
    // read a whole input string from the file
    let contents = fs::read_to_string("puzzle-input.txt").expect("Something went wrong reading the file");
    // split this up into individual rows
    let rows_str: Vec<&str> = contents.split("\n").collect();

    // prepare by creating an appropriately sized array filled with floor cells
    let mut seat_array: [[i32; 92]; 94] = [[-1; 92]; 94];

    // go through the rows and split up that row of text into chars
    for row_index in 0..94 {
        let mut row_chars = rows_str[row_index].chars();
        // for each char that is an 'L', mark the corresponding cell as empty (instead of being a floor tile)
        for col_index in 0..92 {
            if row_chars.next().unwrap() == 'L' {
                seat_array[row_index][col_index] = 0;
            }
        }
    }

    // loops once for every round in simulation
    loop {
        // has the array been changed this round?
        let mut changed_array = false;

        // this is readonly
        // changes from last round are only able to be read at the beginning of the next round
        // think of this like a 'snapshot' of the state of the room LAST round as THIS round is taking place
        // changes are directly commited to seat_array
        let seat_array_this_round = seat_array;

        // go through each pair of indexes
        for row_index in 0..94 {
            for col_index in 0..92 {
                // if this is a floor tile, get out of here!
                // we do not need to do any computation in this case
                if seat_array_this_round[row_index as usize][col_index as usize] == -1 {
                    continue;
                }

                // count occupied seats adjacent to the desired seat
                let mut count_adjacent: i32 = 0;

                // loop through all the real cells (ones which are actually in a valid range - this is checked by get_adjacent()) adjacent to our desired cell
                for adjacent_coords in get_adjacent(row_index, col_index) {
                    // if this cell is a 1 (is a filled seat), count this in count_adjacent
                    if seat_array_this_round[adjacent_coords[0] as usize][adjacent_coords[1] as usize] == 1 {
                        count_adjacent += 1;
                    }
                }

                // we can finally implement the rules!
                // only one rule needs to be checked for a single cell that contains a seat
                // a different rule is applied based on whether this is a filled or unfilled seat
                match seat_array_this_round[row_index as usize][col_index as usize] {
                    0 => {
                        // if there are no adjacent tiles with people, fill this seat!
                        if count_adjacent == 0 {
                            seat_array[row_index as usize][col_index as usize] = 1;
                            // also remember to keep track of this change
                            changed_array = true;
                        }
                    }
                    1 => {
                        // if there are too many adjacent tiles with people, remove the person from this seat!
                        if count_adjacent >= 4 {
                            seat_array[row_index as usize][col_index as usize] = 0;
                            // also remember to keep track of this change
                            changed_array = true;
                        }
                    }
                    // the floor tiles were weeded out earlier
                    _ => panic!()
                }
            }
        }

        // if we managed to get through this round without changing anything, we've reached an equilibrium!
        if !changed_array {
            // count up the number of filled seats in the end (these cells have a value of 1)
            let mut filled_seats = 0;
            for row_index in 0..94 {
                for col_index in 0..92 {
                    if seat_array[row_index][col_index] == 1 {
                        filled_seats += 1;
                    }
                }
            }

            println!("{}", filled_seats);
            break;
        }
    }
}

fn part2() {
    // read a whole input string from the file
    let contents = fs::read_to_string("puzzle-input.txt").expect("Something went wrong reading the file");
    // split this up into individual rows
    let rows_str: Vec<&str> = contents.split("\n").collect();

    // prepare by creating an appropriately sized array filled with floor cells
    let mut seat_array: [[i32; 92]; 94] = [[-1; 92]; 94];

    // go through the rows and split up that row of text into chars
    for row_index in 0..94 {
        let mut row_chars = rows_str[row_index].chars();
        // for each char that is an 'L', mark the corresponding cell as empty (instead of being a floor tile)
        for col_index in 0..92 {
            if row_chars.next().unwrap() == 'L' {
                seat_array[row_index][col_index] = 0;
            }
        }
    }

    // possible directions a person can look
    let directions: [[i32; 2]; 8] = [[1,1], [1,0], [1,-1], [0,1], [0,-1], [-1,1], [-1,0], [-1,-1]];

    // loops once for every round in simulation
    loop {
        // has the array been changed this round?
        let mut changed_array = false;

        // this is readonly
        // changes from last round are only able to be read at the beginning of the next round
        // think of this like a 'snapshot' of the state of the room LAST round as THIS round is taking place
        // changes are directly commited to seat_array
        let seat_array_this_round = seat_array;

        // go through each pair of indexes
        for row_index in 0..94 {
            for col_index in 0..92 {
                // if this is a floor tile, get out of here!
                // we do not need to do any computation in this case
                if seat_array_this_round[row_index as usize][col_index as usize] == -1 {
                    continue;
                }

                // count occupied seats adjacent to the desired seat
                let mut count_adjacent: i32 = 0;

                // here is where part 2 differs from part 1!
                // for every direction a persion can look...
                for direction in &directions {
                    let mut query_index = [row_index, col_index];
                    // ...loop over the indexes going in that direction...
                    loop {
                        // move 1 step in this direction
                        query_index[0] += direction[0];
                        query_index[1] += direction[1];

                        // ...until we encounter...

                        // ...an invalid index (in this case simply break and move onto the next direction), or...
                        if query_index[0]<0 || query_index[0]>93 || query_index[1]<0 || query_index[1]>91 {
                            break;
                        }

                        // ...a cell, (in which case process the value inside)
                        match seat_array_this_round[query_index[0] as usize][query_index[1] as usize] {
                            // if there is a person, take record of this and break
                            1 => {
                                count_adjacent += 1;
                                break;
                            },
                            // if there is no person (but there is a seat), still break because it is not possible to see past this seat
                            0 => break,
                            _ => {}
                        }
                    }
                }

                // we can finally implement the rules!
                // only one rule needs to be checked for a single cell that contains a seat
                // a different rule is applied based on whether this is a filled or unfilled seat
                match seat_array_this_round[row_index as usize][col_index as usize] {
                    0 => {
                        // if there are no adjacent tiles with people, fill this seat!
                        if count_adjacent == 0 {
                            seat_array[row_index as usize][col_index as usize] = 1;
                            // also remember to keep track of this change
                            changed_array = true;
                        }
                    }
                    1 => {
                        // if there are too many adjacent tiles with people, remove the person from this seat!
                        if count_adjacent >= 5 {
                            seat_array[row_index as usize][col_index as usize] = 0;
                            // also remember to keep track of this change
                            changed_array = true;
                        }
                    }
                    // the floor tiles were weeded out earlier
                    _ => panic!()
                }
            }
        }

        // if we managed to get through this round without changing anything, we've reached an equilibrium!
        if !changed_array {
            // count up the number of filled seats in the end (these cells have a value of 1)
            let mut filled_seats = 0;
            for row_index in 0..94 {
                for col_index in 0..92 {
                    if seat_array[row_index][col_index] == 1 {
                        filled_seats += 1;
                    }
                }
            }

            println!("{}", filled_seats);
            break;
        }
    }
}

fn main() {
    part1();
    part2();
}