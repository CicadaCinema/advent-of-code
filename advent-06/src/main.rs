use std::fs;

fn main() {
    // PART 1
    // read a whole string from the input file
    let contents = fs::read_to_string("puzzle-input.txt").expect("Something went wrong reading the file");
    // split big string into pass elements
    let groups: Vec<&str> = contents.split("\n\n").collect();
    // counter for all groups
    let mut count: i32 = 0;

    // iterate over the groups one by one
    for i in 0..groups.len() {
        // new lines don't matter here
        let group = groups[i as usize].replace("\n", "");
        // unique characters in this group
        let mut chars_used: Vec<char> = vec![];

        // if character is not used yet, add it to the vector above
        for character in group.chars() {
            if !chars_used.contains(&character) {
                chars_used.push(character);
            }
        }

        // increment global count by count in this group
        count += chars_used.len() as i32;
    }

    println!("{}", count);

    // PART 2
    // reset count
    count = 0;

    // iterate over the groups one by one
    for i in 0..groups.len() {
        // people do matter here
        let people: Vec<&str> = groups[i as usize].split("\n").collect();
        // assume everyone chose the same chars as person 1 - then remove chars from this vector as needed
        let mut common_chars: Vec<char> = people[0].chars().collect();

        // iterate over the people in this group
        for person in people {
            // we cannot remove elements from a vector while iterating over it - so store the indexes to be removed and remove them all later
            let mut common_chars_scheduled_del: Vec<i32> = vec![];
            // find the characters this person picked
            let person_chars: Vec<char> = person.chars().collect();

            // for every common character (thus far), verify its presence in this person's list of character
            // if it is not present, add its index to a list scheduled for deletion
            for j in 0..common_chars.len() {
                if !person_chars.contains(&common_chars[j]) {
                    common_chars_scheduled_del.push(j as i32);
                }
            }

            // sort these indexes in descending order - this means that any small indexes are not affected by the removal of elements at larger indexes
            common_chars_scheduled_del.sort();
            common_chars_scheduled_del.reverse();

            // go ahead and remove these indexes from the common chars vector
            for index in common_chars_scheduled_del {
                common_chars.remove(index as usize);
            }
        }

        // increment global count by count in this group
        count += common_chars.len() as i32;
    }

    println!("{}", count);
}
