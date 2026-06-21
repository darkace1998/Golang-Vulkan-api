import re

with open("command_test.go", "r") as f:
    text = f.read()

text = re.sub(r'func BenchmarkCommandBufferBatchAdd.*?// TestCreateThreadLocalCommandPool', '// TestCreateThreadLocalCommandPool', text, flags=re.DOTALL)

with open("command_test.go", "w") as f:
    f.write(text)
