fn main() {
    // PART 1
    // goal and bus ids
    let goal = 1002632;
    let buses: [i32; 9] = [23,41,829,13,17,29,677,37,19];

    // stores current best waiting time and the corresponding id
    let mut answer: (i32, i32) = (99999, 0);

    for bus_index in 0..9 {
        // for each bus, calculate the required waiting time using some integer division and math
        let waiting_time: i32 = (buses[bus_index] * (1 + (goal / buses[bus_index]))) - goal;
        // if this is a new best in terms of waiting time, save it to answer
        if waiting_time < answer.0 {
            answer.0 = waiting_time;
            answer.1 = buses[bus_index];
        }
    }

    // show product
    println!("{}", answer.0 * answer.1);

    // PART 2
    let buses_contest = [23,0,0,0,0,0,0,0,0,0,0,0,0,41,0,0,0,0,0,0,0,0,0,829,0,0,0,0,0,0,0,0,0,0,0,0,13,17,0,0,0,0,0,0,0,0,0,0,0,0,0,0,29,0,677,0,0,0,0,0,37,0,0,0,0,0,0,0,0,0,0,0,0,19];
    let mut buses_delay: [i32; 9] = [0; 9];
    for index in 1..9 {
        buses_delay[index] = buses_contest.iter().position(|r|r==&buses[index]).unwrap() as i32;
    }
    
    // buses  [23, 41, 829, 13, 17, 29, 677, 37, 19]
    // delays [0,  13, 23,  36, 37, 52, 54,  60, 73]

    let mut check_buses = vec![buses[0]];
    let mut check_delays = vec![buses_delay[0]];

    let mut memory: i64 = 0;
    let mut current_try: i64 = 0;
    let mut current_increment: i64 = buses[0] as i64;

    for i in 1..(buses.len()) {
        check_buses.push(buses[i]);
        check_delays.push(buses_delay[i]);

        'outer: loop {
            //println!("{} {}", current_try, current_increment);
            current_try += current_increment;
            for check_index in 0..(i+1) {
                if (current_try + (check_delays[check_index] as i64)) % (check_buses[check_index] as i64) != 0 {
                    continue 'outer;
                }
            }

            if memory == 0 {
                memory = current_try;
            } else {
                current_increment = current_try - memory;
                current_try = memory;
                memory = 0;
                break;
            }
        }
    }

    println!("{}", current_try);
}
