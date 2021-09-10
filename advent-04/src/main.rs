use std::fs;

fn main() {
    // PART 1
    // read a whole string from the input file
    let contents = fs::read_to_string("puzzle-input.txt").expect("Something went wrong reading the file");
    // split big string into pass elements
    let passes: Vec<&str> = contents.split("\n\n").collect();

    // keep track of fields and prepare to count valid passes
    let required_fields:[&str; 8] = ["byr", "iyr", "eyr", "hgt", "hcl", "ecl", "pid", "cid"];
    let mut valid_passes_count: i32 = 0;

    // iterate over the passes one by one
    for i in 0..passes.len() {
        // treat a new line the same way as a space and then split the pass into a vector
        let element = passes[i as usize].replace("\n", " ");
        let one_pass: Vec<&str> = element.split(" ").collect();
        // prepare to validate every field (exclusing cid)
        let mut present_fields = [false; 7];

        // iterate over the fields in this pass
        for j in one_pass {
            // the first three chars determine the field name
            let field_name = &j[..3];
            if field_name != "cid" {
                // if this field is not cid, record it as present by making the corresponding index true in present_fields
                let field_index = required_fields.iter().position(|&r| r == field_name).unwrap();
                present_fields[field_index as usize] = true;
            }
        }

        // if all 7 fields are true (in present_fields), increment the count of valid passes
        if present_fields.iter().filter(|&n| *n == true).count() == 7 {
            valid_passes_count += 1;
        }
    }

    println!("{} valid passes", valid_passes_count);

    // PART 2
    // new variable to count the number of valid passes
    let mut valid_passes_count: i32 = 0;

    // iterate over the passes one by one
    for i in 0..passes.len() {
        // treat a new line the same way as a space and then split the pass into a vector
        let element = passes[i as usize].replace("\n", " ");
        let one_pass: Vec<&str> = element.split(" ").collect();
        // prepare to validate every field (exclusing cid)
        let mut present_fields = [false; 7];

        // iterate over the fields in this pass
        for j in one_pass {
            // the first three chars determine the field name
            let field_name = &j[..3];
            // chars from five to the last determine the field value
            let field_value = &j[4..];

            if field_name != "cid" {
                // if this field is not cid, record it as present by making the corresponding index true in present_fields
                let field_index = required_fields.iter().position(|&r| r == field_name).unwrap();
                
                // match the field name
                present_fields[field_index as usize] = match field_name {
                    // first three field types are just numbers - they must be in a specified range to be valid
                    "byr" => {
                        let num: i32 = field_value.parse().unwrap();
                        num >= 1920 && num <= 2002
                    },
                    "iyr" => {
                        let num: i32 = field_value.parse().unwrap();
                        num >= 2010 && num <= 2020
                    },
                    "eyr" => {
                        let num: i32 = field_value.parse().unwrap();
                        num >= 2020 && num <= 2030
                    },
                    // determine unit first (last two chars) and then apply corresponding range
                    "hgt" => {
                        let unit = &field_value[(field_value.len() - 2)..];
                        let real_value = &field_value[..(field_value.len() - 2)];
                        match unit {
                            "cm" => {
                                let num: i32 = real_value.parse().unwrap();
                                num >= 150 && num <= 193
                            },
                            "in" => {
                                let num: i32 = real_value.parse().unwrap();
                                num >= 59 && num <= 76
                            },
                            _ => false,
                        }
                    },
                    // ensure that the length is 7 and that the first character is a #
                    // then go through the characters one by one and ensure they are either numeric or a lowercase letter
                    // valid is true by default, but make it false as soon as one of these conditions isn't fulfilled
                    "hcl" => {
                        let mut valid = true;
                        if field_value.len()==7 && field_value.chars().nth(0).unwrap()=='#' {
                            for character_index in 1..7 {
                                let character = field_value.chars().nth(character_index).unwrap();
                                if (character.is_numeric() || (character.is_alphabetic() && character.is_lowercase())) != true {
                                    valid = false;
                                }
                            }
                        } else {
                            valid = false;
                        }
                        valid
                    },
                    // this field is simple to validate - just see if the value is in the list of allowed ones
                    "ecl" => {
                        ["amb", "blu", "brn", "gry", "grn", "hzl", "oth"].contains(&field_value)
                    },
                    // essentially the same thing as hcl, but we only check if there are 9 chracters and each of them is a number
                    "pid" => {
                        let mut valid = true;
                        if field_value.len()==9 {
                            for character in field_value.chars() {
                                if character.is_numeric() != true {
                                    valid = false;
                                }
                            }
                        } else {
                            valid = false;
                        }
                        valid
                    },
                    // panic if we see a field that is not accounted for
                    _ => panic!(),
                }
            }
        }

        // if all 7 fields are true (in present_fields), increment the count of valid passes
        if present_fields.iter().filter(|&n| *n == true).count() == 7 {
            valid_passes_count += 1;
        }
    }

    println!("{} valid passes", valid_passes_count);
}
