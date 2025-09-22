// Be sure to be in cd \javascript
// Run "npm run test" in terminal
// Resource : https://cp-algorithms.com/algebra/sieve-of-eratosthenes.html
class Sieve {
    segSieve(s) {
        if (s < 0) return null;

        let limit = Math.ceil(s * (Math.log(s) + Math.log(Math.log(s)) + 2));
        let segmentSize = 1e6;                                          // 1 million

        let sqrtLimit = Math.floor(Math.sqrt(limit));
        let sqrtPrimes = new Uint8Array(sqrtLimit + 1).fill(true);
        sqrtPrimes[0] = sqrtPrimes[1] = false;

        let basePrimes = [];
        for (let i = 2; i <= sqrtLimit; i++) {
            if (sqrtPrimes[i]) {
                basePrimes.push(i);
                for (let j = i * i; j <= sqrtLimit; j += i) {
                    sqrtPrimes[j] = false;
                }
            }
        }

        let count = 0;
        for (let low = 2; low <= limit; low += segmentSize) {           // low is always the new starting point of the buffer on each iteration
            let high = Math.min(low + segmentSize - 1, limit);          // set high, with a cap to the calc of limit so it doesnt exceed bounds
            let mark = new Uint8Array(high - low + 1).fill(true);       // faster than reseting buffer

            for (let p of basePrimes) {
                let start = Math.max(p * p, Math.ceil(low / p) * p);    // first multiple from prime
                for (let j = start; j <= high; j += p) {                // checks from first multiple to each basePrime on each iteration
                    mark[j - low] = false;                              // set to the 0th of mark[]
                }
            }
            for (let i = low; i <= high; i++) {                         // iterate throw the pos/neg truth flags in mark[]
                if (mark[i - low]) {
                    if (count === s)
                    return i;
                    count++;
                }
                
            }
        }
        throw new Error ("Error finding Segmented Prime");
    }

    NthPrime(n) {                                                       // Main API call
        // if number is larger than 1000000 call Segmented function
        if (n >= 10000000) {
            return this.segSieve(n);
        }
        if (n < 0) {
            n = 0;
        }

        let uBound =
            n < 10
                ? 30
                : Math.ceil(n * (Math.log(n) + Math.log(Math.log(n)) + 2));

        let validatePrimes = new Uint8Array(uBound + 1).fill(true);

        validatePrimes[0] = false;
        validatePrimes[1] = false;

        // sieve of eratosthenes
        for (let i = 2; i * i <= uBound; i++) {
            if (validatePrimes[i]) {
                for (let j = i * i; j <= uBound; j += i) {
                    validatePrimes[j] = false;
                }
            }
        }

        // get nth prime
        let primeCount = 0;
        for (let i = 2; i <= uBound; i++) {
            if (validatePrimes[i]) {
                if (primeCount === n) {
                    return i;
                }
                primeCount++;
            }
        }

    }
}

module.exports = Sieve;
