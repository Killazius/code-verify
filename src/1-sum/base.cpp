#include <iostream>

int sum(int, int);

int main(int argc, char* argv[]) {
    if (argc != 3) {
        std::cerr << "Usage: " << argv[0] << " <num1> <num2>" << std::endl;
        return 1;
    }

    int num1 = std::stoi(argv[1]);
    int num2 = std::stoi(argv[2]);

    int result = sum(num1, num2);
    std::cout << result << std::endl;

    return 0;
}
//
