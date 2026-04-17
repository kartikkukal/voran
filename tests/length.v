int length(char* string) {
    int i = 0;
    char c = string[i];

    while (c != '\0') {
        i += 1;
        c = string[i];
    }

    return i;
}

int main() {
    char[] hello = "Unknown String!";
    print(length(hello));
}