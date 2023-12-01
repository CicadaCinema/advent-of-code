use std::fs;

fn part1() {
    // read a whole input string from the file
    let contents = fs::read_to_string("puzzle-input.txt").expect("Something went wrong reading the file");
    // split this up into individual rules
    let rules_raw: Vec<&str> = contents.split("\n").collect();

    // this contains all we need to know about the rules
    // the outer bag is the first element, the inner bags are the second element
    let mut rules: Vec<(&str, Vec<String>)> = vec![];

    // process rules one by one
    for rule in rules_raw.iter() {
        // fint the outer bag in the string
        let outer_bag_index = rule.find(" bags contain").unwrap();
        let outer_bag = &rule[..outer_bag_index];

        // create data structure for storing this rule
        let mut this_rule: (&str, Vec<String>) = (outer_bag, vec![]);

        // try t0 separate out the bag colours, but remove the first and last element because they are not useful
        let mut string_split: Vec<&str> = rule.split(" bag").collect();
        string_split.remove(0);
        string_split.pop();

        // go over each inner bag one by one
        for inner_bag in string_split.iter() {
            // this is chaos, I don't know how this works - I had to wrestle quite a bit with the borrow checker
            // the function of this code is to give me the colour of the bag - I know that colour is expressed by two words
            let inner_bag_split: Vec<&str> = inner_bag.split(" ").collect();
            let mut this_inner_bag: String = String::new();
            this_inner_bag.push_str(&*inner_bag_split.iter().nth(inner_bag_split.len()-2).unwrap().to_owned().to_string());
            this_inner_bag.push_str(&" ".to_owned());
            this_inner_bag.push_str(&*inner_bag_split.iter().nth(inner_bag_split.len()-1).unwrap().to_owned().to_string());

            // store this colour
            this_rule.1.push(String::from(&this_inner_bag));
        }

        // store this rule
        rules.push(this_rule);
    }

    let mut possible_outer_bags: Vec<&str> = vec![];
    let mut inner_bags_to_check: Vec<&str> = vec![&"shiny gold"];

    // I have no clue how this is supposed to work...
    loop {
        for i in (0..inner_bags_to_check.len()).rev() {
            let this_inner_bag_check = inner_bags_to_check[i];

            for rule in rules.iter() {
                if rule.1.iter().find(|&r| r==this_inner_bag_check).is_some() {
                    if !possible_outer_bags.iter().find(|&r| r==&rule.0).is_some() {
                        possible_outer_bags.push(rule.0);
                    }
                    inner_bags_to_check.push(rule.0);
                }
            }

            inner_bags_to_check.remove(i);
        }
        
        if inner_bags_to_check.len() == 0 {
            break;
        }
    }

    println!("{:?}", possible_outer_bags.len());
}

fn part2() {
    // read a whole input string from the file
    let contents = fs::read_to_string("puzzle-input.txt").expect("Something went wrong reading the file");
    // split this up into individual rules
    let rules_raw: Vec<&str> = contents.split("\n").collect();

    // this contains all we need to know about the rules
    // the outer bag is the first element, the inner bags are the second element
    let mut rules: Vec<(&str, Vec<(String, i32)>)> = vec![];

    // process rules one by one
    for rule in rules_raw.iter() {
        // fint the outer bag in the string
        let outer_bag_index = rule.find(" bags contain").unwrap();
        let outer_bag = &rule[..outer_bag_index];

        // create data structure for storing this rule
        let mut this_rule: (&str, Vec<(String, i32)>) = (outer_bag, vec![]);

        // try t0 separate out the bag colours, but remove the first and last element because they are not useful
        let mut string_split: Vec<&str> = rule.split(" bag").collect();
        string_split.remove(0);
        string_split.pop();

        // go over each inner bag one by one
        for inner_bag in string_split.iter() {
            // this is chaos, I don't know how this works - I had to wrestle quite a bit with the borrow checker
            // the function of this code is to give me the colour of the bag - I know that colour is expressed by two words
            let inner_bag_split: Vec<&str> = inner_bag.split(" ").collect();
            let mut this_inner_bag: String = String::new();

            // this finds the str representing the number of (inner) bags of this kind
            let number_of_bags_str: &str = &*inner_bag_split.iter().nth(inner_bag_split.len()-3).unwrap().to_owned().to_string();

            // actually the line above doesn't always work, because some outer bags do not allow any number of inner bags at all!
            // do not do anything - let the rules vector remain empty
            if number_of_bags_str == "contain" {
                // -> rightly so, because this bag contain no other bags
            } else {
                // now we know that there is some number in number_of_bags_str, so we can safely convert it to an i32
                let number_of_bags_int: i32 = number_of_bags_str.parse::<i32>().unwrap();

                this_inner_bag.push_str(&*inner_bag_split.iter().nth(inner_bag_split.len()-2).unwrap().to_owned().to_string());
                this_inner_bag.push_str(&" ".to_owned());
                this_inner_bag.push_str(&*inner_bag_split.iter().nth(inner_bag_split.len()-1).unwrap().to_owned().to_string());
    
                // store this colour, and the corresponding number of bags
                this_rule.1.push((String::from(&this_inner_bag), number_of_bags_int));
            }
        }

        // store this rule
        rules.push(this_rule);
    }

    println!("{:?}", find_bags(&rules, &"shiny gold"));
}

fn find_bags(all_rules: &Vec<(&str, Vec<(String, i32)>)>, bag_colour: &str) -> i32 {
    // search for the relevant rule (the bag we are currently on has to be the outer bag)
    let relevant_rule: &(&str, Vec<(String, i32)>) = all_rules.iter().find(|&r| r.0==bag_colour).unwrap();
    // begin counting the contents of this bag
    let mut bag_count: i32 = 0;

    if relevant_rule.1.len() == 0 {
        // if this bag can contain no other bags, return 0 straight away!
        return 0;
    }

    for i in 0..relevant_rule.1.len() {
        // for every separate inner bag colour, we have to count:
        // * the number of inner bags of that colour
        // * the number of bags contained within this colour inner bag - so the number of bags of this colour MULTIPLIED BY the number of bags each one contains! (by recursion)
        bag_count += relevant_rule.1[i].1;
        bag_count += relevant_rule.1[i].1 * find_bags(all_rules, &relevant_rule.1[i].0);
    }

    return bag_count;
}

fn main() {
    part1();
    part2();
}