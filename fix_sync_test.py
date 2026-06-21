with open("synchronization_test.go", "r") as f:
    lines = f.readlines()

with open("synchronization_test.go", "w") as f:
    for line in lines[:149]:
        f.write(line)
