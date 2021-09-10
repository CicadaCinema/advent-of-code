fn main() {
    // PART 1
    let mut adapters = vec![152,18,146,22,28,133,114,67,19,37,66,14,90,163,26,149,71,106,46,143,145,12,151,105,58,130,93,49,74,83,129,122,63,134,86,136,166,169,159,3,178,88,103,97,110,53,125,128,9,15,78,1,50,87,56,89,60,139,113,43,36,118,170,96,135,23,144,153,150,142,95,180,35,179,80,13,115,2,171,32,70,6,72,119,29,79,27,47,107,73,162,172,57,40,48,100,64,59,175,104,156,94,77,65];
    let mut count_1 = 0;
    let mut count_3 = 0;

    // sort and also don't forget the last adapter, which has a value of 3 greater than the next-greatest adapter
    adapters.sort();
    adapters.push(adapters.last().unwrap() + 3);

    // wall outlet has a value of 0
    let mut prev = 0;
    // go through the adapters one by one
    for adapter in &adapters {
        // if the difference is 1 or 3, count that up
        match adapter - prev {
            1 => count_1 += 1,
            3 => count_3 += 1,
            _ => {}
        }
        // remember this adapter as previous for the next cycle
        prev = *adapter;
    }

    // print the product
    println!("{}", count_1 * count_3);


    // PART 2

    // Theory:
    // when we compute the difference between adjacent terms in the (sorted) adapter array, we will find diffs of 1, 2 and 3
    // our goal is to find the number of ways to concatenate these diffs into a sequence which still has elements 1, 2 and 3
    // the code below helps us to find sequences which can be concatenated (these can consist of 2s and 1s, but no 3s)

    // wall outlet has a value of 0
    let mut prev = 0;
    // vector used to store the current sequence of possible characters
    let mut seq: Vec<i32> = vec![];
    for adapter in &adapters {
        // get difference between this element and previous
        let diff = adapter-prev;

        // discovery 1: this is never run - there are no 2s at all!
        if diff == 2 {
            println!("there is a 2!");
        }

        match diff {
            // if diff is 1 or 2, push it to the vector
            1 => seq.push(diff),
            2 => seq.push(diff),
            // otherwise the sequence of concatenateable characters terminates, so print out this seq and start afresh
            _ => {
                if seq.len() > 0 {
                    println!("{:?}", seq);
                    seq = vec![];
                }
            }
        }

        // remember this adapter as previous for the next cycle
        prev = *adapter;
    }
    // make sure that the last substring is displayed, if it exists
    if seq.len() > 0 {
        println!("{:?}", seq);
    }

    // thanks to discovery 1, the problem is now much easier
    // it is not possible to concat a single diff, so we are left with four [1, 1]s, six [1, 1, 1]s and nine [1, 1, 1, 1]s

    // concat options are as follows:
    // [1, 1] or [2]
    // [1, 1, 1] or [1, 2] or [2, 1] or [3]
    // [1, 1, 1, 1] or [1, 1, 2] or [1, 2, 1] or [2, 1, 1] or [2, 2] or [1, 3] or [3, 1]

    // therefore the answer is just
    // 2^4 * 4^6 * 7^9
    // when plugged into wolfram alpha, this computes to 2644613988352
}
