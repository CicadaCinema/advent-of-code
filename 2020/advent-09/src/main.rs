use std::fs;

// checks if the number at the specified index is valid
// that is to say, whether any pair of the last 25 elements can sum to the number at index
fn check_index(number_vec: &Vec<i64>, index: usize) -> bool {
    let goal_num = number_vec[index];
    let range = &number_vec[(index-25_usize)..index];

    for i in range {
        if range.contains(&(goal_num-i)) {
            return true;
        }
    }

    return false;
}

fn main() {
    // PART 1
    // read a whole input string from the file
    let contents = fs::read_to_string("puzzle-input.txt").expect("Something went wrong reading the file");
    // split this up into individual numbers
    let number_input: Vec<&str> = contents.split("\n").collect();
    // collect these in a vector
    let mut number_input_int: Vec<i64> = Vec::new();
    for i in number_input {
        number_input_int.push(i.parse().unwrap());
    }

    // required storage for part 2
    let mut invalid_number: i64 = 0;

    // check every index from the 26th element to the end of the vector
    for probe_index in 25..(number_input_int.len()) {
        // if it is invalid, hooray! we can now record this special invalid number and break out from the loop
        if !check_index(&number_input_int, probe_index as usize) {
            invalid_number = number_input_int[probe_index];
            break;
        }
    }

    println!("{}", invalid_number);

    // PART 2

    // this uses my "inchworm" algorithm to find a contiguous set which sums to the goal - invalid_number

    // at first, take the test set to consist of the first two elements and therefore calculate the sum from them
    let mut current_index_begin: usize = 0;
    let mut current_index_end: usize = 1;
    let mut current_sum: i64 = number_input_int[current_index_begin] + number_input_int[current_index_end];

    // loop until we find a valid contiguous set (a runtime error occurs eventually, if there is no such set)
    loop {
        if current_sum < invalid_number {
            // we can progress forwards by incrementing the end index (and then add that number to the sum)
            current_index_end += 1;
            // there will be a runtime error here if the problem of finding a satisfactory contiguous set is unsolveable
            // this is because we will reach the end of the number_input_int vector without finding the required sum, so we will try to index
            // the vector with current_index_end, which will be too high (in fact it will be equal to the length of the vector)
            current_sum += number_input_int[current_index_end];
        } else if current_sum > invalid_number {
            // we need to reduce the sum by incrementing the start index (only after decrementing the sum by this amount)
            current_sum -= number_input_int[current_index_begin];
            current_index_begin += 1;
        } else if current_sum == invalid_number {
            // solved! we found a matching contiguous set
            break;
        }
    }

    // now we know our range, so we must find the min and max number within this range

    // assume the first number is both the max and the min - for now
    let mut min_number: i64 = number_input_int[current_index_begin];
    let mut max_number: i64 = number_input_int[current_index_begin];

    for element_index in (current_index_begin + 1)..(current_index_end + 1) {
        // cycle through every element in the range (but we don't need to check the first one)
        let current_element: i64 = number_input_int[element_index];
        if current_element < min_number {
            // if the current element is smaller than the min, make it the min
            min_number = current_element;
        } else if current_element > max_number {
            // if the current element is greater than the max, make it the max
            max_number = current_element;
        }
    }

    // the answer to part 2 is the sum of these two numbers
    println!("{}", max_number + min_number);
}
