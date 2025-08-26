print("=====")
local startTime = os.clock()
local total = 0
for i = 1, 1000000 do
    total = total + i
end
local endTime = os.clock()
local elapsedMs = (endTime - startTime) * 1000
print(string.format("Time taken: %.3f milliseconds", elapsedMs))
print(string.format("Final result: %d", total))