with open("map.pbm", "r") as f:
    f.readline() # Skip magic number
    f.readline() # Skip GIMP comment
    width, height = [int(x) for x in f.readline().split(" ")]

    print(" "*4, end="")

    printed = 0
    row = 0
    col = 0
    for val in f.read().replace("\n", ""):
        if col == width-1:
            col = 0
            row += 1
        else:
            col += 1

        if val == "1":
            print("{"+str(col).rjust(2)+", "+str(row).rjust(2)+"}, ", end="")
            printed += 1
            if printed % 8 == 0:
                print("\n"+" "*4, end="")

print()
