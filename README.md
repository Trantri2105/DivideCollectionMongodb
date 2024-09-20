# Chia Collection trong MongoDB

## Giới thiệu
- Project chia 1 collection trong mongodb thành nhiều collection theo các giá trị trong thuộc tính `entity_type`

## Flow

1. Kết nối với Mongodb
2. Lấy dữ liệu trong collection cần chia và phân chia các dữ liệu này vào các collection khác

## Set up
- Để code hoạt động, ta cần phải khai báo giá trị cho một số biến sau đây
  - `mongodbURI`
  - `database`
  - `currentCollection`: tên collection cần chia
  - `batchSize`: số document lấy từ database trong 1 đợt và là số document insert vào database trong 1 đợt
  - `newCollection`: map chứa giá trị `entity_type` và collection tương ứng. VD:
    ```
    newCollection["device"] = "device_attribute"
    newCollection["user"] = "user_attribute"
    ```
    Các document có giá trị `entity_type` là `device` sẽ được lưu vào `device_attribute` còn `user` sẽ được lưu vào `user_attribute`

## Test
1. Tạo 2000 record và lưu vào trong collection `attribute`, 1000 record có `entity_type` là `user` và 1000 record là `device`. Mỗi record có dạng
   ```
   {
      entity_id:  ,
      entity_type:  
   }
   ```
   hoặc
   ```
   {
        entity_id: "Entity_1",
        entity_type: "device",
        name: "abcd",
        ab: 1,
		state: true,
        abcd: {
			a : "bc",
			num: 123
		}
   }
   ```
   
2. Lấy record đã lưu trong attribute và chia vào 2 slice `userDocuments` và `deviceDocuments`
3. Chạy hàm main để tiến hành chia collection
4. Lấy các record trong `user_attribute` và so sánh với `userDocuments`
5. Lấy các record trong `device_attribute` và so sánh với `deviceDocuments`