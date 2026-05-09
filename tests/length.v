i32 number
number = 27

i32 max
max = 0

while number != 1 {
    if number % 2 != 0 {
        number = number * 3 + 1
    }
    if max < number {
        max = number
    }
    if number % 2 == 0 {
        number = number / 2
    }
}