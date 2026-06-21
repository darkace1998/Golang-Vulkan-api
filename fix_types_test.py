with open("types_test.go", "r") as f:
    text = f.read()

text = text.replace("func TestResultHelpers(", "func TestResultHelpersTypes(")

with open("types_test.go", "w") as f:
    f.write(text)
