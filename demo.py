import time

print(f"=====")
start_time = time.time()
total = 0
for i in range(1, 1000001):
    total += i
end_time = time.time()
elapsed_ms = (end_time - start_time) * 1000
print(f"Time taken: {elapsed_ms:.3f} milliseconds")
print(f"Final result: {total}")


print(f"=====")
start_time = time.time()
total = 0
for i in range(1, 1000001):
    total += i
end_time = time.time()
elapsed_ms = (end_time - start_time) * 1000
print(f"Time taken: {elapsed_ms:.3f} milliseconds")
print(f"Final result: {total}")

print(f"=====")
start_time = time.time()
total = 0
for a in range(1, 1001):
    for b in range(1, 1001):
        total += a + b
    end_time = time.time()
end_time = time.time()

elapsed_ms = (end_time - start_time) * 1000
print(f"Time taken: {elapsed_ms:.3f} milliseconds")
print(f"Final result: {total}")