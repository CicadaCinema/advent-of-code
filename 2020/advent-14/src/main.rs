use std::fs;

fn to_binary_array(input_num: i64) -> [u8; 36] {
    let mut current_num = input_num;
    let mut output = [0_u8; 36];

    for out_index in (0..36).rev() {
        if current_num%2 == 1 {
            output[out_index as usize] = 1;
            current_num -= 1;
        }
        current_num = current_num / 2;
    }

    output
}

fn set_mem(input_array: [u8; 36], mask: [u8; 36]) -> [u8; 36] {
    let mut output_array = input_array;

    for bit_index in 0..36 {
        if mask[bit_index] != 2 {
            output_array[bit_index] = mask[bit_index];
        }
    }

    output_array
}

fn get_mem_size(memory: [[u8; 36]; 100000]) -> u64 {
    let mut sum = 0;

    for number in memory.iter() {
        for bit_index in 0..36 {
            let power: u64 = 2_u64.pow(35-bit_index);
            sum += power*(number[bit_index as usize] as u64);
        }
    }

    sum
}

// technically the mask here is the COMBINED memory address, not just the mask
fn get_possible_mem_addresses(mask: [u8; 36]) -> Vec<[u8; 36]> {
    let mut output_vector: Vec<[u8; 36]> = vec![mask];

    'outer: loop {
        for i in (0..(output_vector.len())).rev() {
            let this_address = output_vector[i];
            for bit_index in 0..36 {
                if this_address[bit_index] == 2 {
                    let mut candidate_1 = this_address;
                    let mut candidate_2 = this_address;
                    candidate_1[bit_index] = 0;
                    candidate_2[bit_index] = 1;
                    output_vector.push(candidate_1);
                    output_vector.push(candidate_2);
                    output_vector.remove(i);
                    continue 'outer;
                }
            }
        }
        break;
    }

    output_vector
}

fn main() {
    // PART 1

    // read a whole input string from the file
    let contents = fs::read_to_string("puzzle-input.txt").expect("Something went wrong reading the file");
    // split this up into individual instructions
    let instructions_read: Vec<&str> = contents.split("\n").collect();

    let mut current_mask = [2_u8; 36];


    // arbitrary size - this only works because the addresses in the test case are small
    let mut memory = [[0_u8; 36]; 100000];


    for i in instructions_read.iter() {
        let split: Vec<&str> = i.split("_").collect();
        match split.len() {
            // change the mask
            1 => {
                let mut mask_chars = split[0].chars();
                for bit_index in 0..36 {
                    current_mask[bit_index] = mask_chars.next().unwrap().to_digit(10).unwrap() as u8;
                }
            }
            // set value to memory
            2 => {
                memory[split[0].parse::<usize>().unwrap()] = set_mem(to_binary_array(split[1].parse().unwrap()), current_mask);
            }
            _ => panic!()
        }
    }

    
    println!("{}", get_mem_size(memory));


    // PART 2
    let mut memory_2: Vec<([u8; 36], u64)> = vec![];

    for i in instructions_read.iter() {
        let split: Vec<&str> = i.split("_").collect();
        match split.len() {
            // change the mask
            1 => {
                let mut mask_chars = split[0].chars();
                for bit_index in 0..36 {
                    current_mask[bit_index] = mask_chars.next().unwrap().to_digit(10).unwrap() as u8;
                }
            }
            // set value to mask
            2 => {
                let mut address = to_binary_array(split[0].parse().unwrap());
                // combine mask
                for bit_index in 0..36 {
                    if current_mask[bit_index] != 0 {
                        address[bit_index] = current_mask[bit_index];
                    }
                }
                // 'write' value to all these memory addresses
                for possible_address in get_possible_mem_addresses(address) {
                    memory_2.push((possible_address, split[1].parse().unwrap()));
                }
            }
            _ => panic!()
        }
    }

    // now sum up the results, making not sure not to 'read' from the same memory location twice

    let mut read_addresses: Vec<[u8; 36]> = vec![];
    let mut sum: u64 = 0;
    for pseudo_read_index in (0..(memory_2.len())).rev() {
        println!("{}", pseudo_read_index);
        if !read_addresses.contains(&memory_2[pseudo_read_index].0) {
            sum += memory_2[pseudo_read_index].1;
            read_addresses.push(memory_2[pseudo_read_index].0);
        }
    }

    println!("{}", sum);
}