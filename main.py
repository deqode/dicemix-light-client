p = 2**61 - 1

def p1():
  my_id = 1
  my_msg_hash = 10

  peers = [[40,2], [60,3]]
  my_msg_num = 1
  total_msg_num = 3

 

  peers[0].append([1518585883438452078, 738716743465821554, 738716743465821554])
  peers[1].append([185161639047148264, 2230111456717881291, 2230111456717881291])

  dc_combined = [1775709972506440409, 1642857818243685069, 1642857818243685069]
  for peer in peers:
    for i in range(total_msg_num):
      dc_combined[i] = reduce(dc_combined[i] + peer[2][i])
  print("P1 DC-COMBINED[] = ", dc_combined)


def p2():
  my_id = 2
  my_msg_hash = 20

  peers = [[40,1], [80,3]]
  my_msg_num = 1
  total_msg_num = 3

  my_dc = [0,0,0]

  for j in range(my_msg_num):
    for i in range(total_msg_num):
      my_dc[i] = reduce(my_dc[i] + (my_msg_hash ** (i+1)))

  for peer in peers:
    for i in range(total_msg_num):
      if my_id < peer[1]:
        my_dc[i] = reduce(my_dc[i] + (p - peer[0]))
      else:
        my_dc[i] = reduce(my_dc[i] + peer[0])
  print("P2 DC[] = ", my_dc)

  peers[0].append([2305843009213693861, 0, 900])
  peers[1].append([170, 1040, 27140])

  dc_combined = my_dc
  for peer in peers:
    for i in range(total_msg_num):
      dc_combined[i] = reduce(dc_combined[i] + peer[2][i])
  print("P2 DC-COMBINED[] = ", my_dc)

def p3():
  my_id = 3
  my_msg_hash = 30

  peers = [[60,1], [80,2]]
  my_msg_num = 1
  total_msg_num = 3

  my_dc = [0,0,0]

  for j in range(my_msg_num):
    for i in range(total_msg_num):
      my_dc[i] = reduce(my_dc[i] + (my_msg_hash ** (i+1)))

  for peer in peers:
    for i in range(total_msg_num):
      if my_id < peer[1]:
        my_dc[i] = reduce(my_dc[i] + (p - peer[0]))
      else:
        my_dc[i] = reduce(my_dc[i] + peer[0])
  print("P3 DC[] = ", my_dc)

  peers[0].append([2305843009213693861, 0, 900])
  peers[1].append([2305843009213693931, 360, 7960])

  dc_combined = my_dc
  for peer in peers:
    for i in range(total_msg_num):
      dc_combined[i] = reduce(dc_combined[i] + peer[2][i])
  print("P3 DC-COMBINED[] = ", my_dc)

def reduce(x):
  return x % p


p1()

