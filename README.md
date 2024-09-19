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