# Sorted string table

## Components

1. Memtable

- Hold the in memory sorted structure
- Currently use AVL tree

2. Disk segment files

- Currently support 1 layer
- The format is (default with md5 hash)
  Key size | key | value size | value
  -|-|-|-
  8 bytes | | 8 bytes |

- Key size, value size can be config with external value or infer from hashing value

3. Disk index files

- Hold key and offset in segment file, since it is sorted, it can be binary searched
- Format (default with md5 hash)
  Key | Offset
  -|-
  16 bytes | 8 bytes
- TODO: Add sampling to hold only 1 every N keys

## Procedure

1. Disk flushing: based on size only

2. Searching

- Key is hashed then sorted. It may not work well with range query. If you want to use range query, then set the sort option as false
- Memtable first
- Then using index files
  - Binary search the key, then using the pos to get the value from segment file
- Does not load the value into memtable on search
