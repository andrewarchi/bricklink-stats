# BrickLink Order Counts

BrickLink, a hobbyist LEGO marketplace, recently surpassed 10 million orders and
I was the 9999999th and 10000000th. In order to do this, I wrote this simple
script to monitor the current order count and predict when 10 million would be
reached.

![Order #10000000](order_10000000.png)
![Order #9999999](order_9999999.png)

BrickLink provides an API, but has no way to determine the current number of
orders placed. The order detail page
(`https://www.bricklink.com/orderDetail.asp?ID=9999999`) can be used to
indirectly find this. When the order is not mine, I get a 403 not authorized
page and invalid order IDs give order not found pages. These redirects indicate
if an order exists.

The order count is estimated with a linear regression from my previous orders,
then the accurate count is found by a binary search. As further orders are made,
they are detected and reported.

Output columns:
- Order date/time
- Order ID
- Time since previous order
- Average time since previous order
- Estimated target time
- Estimated time until target

Below is the output for when 10 million was reached. It is notable that the
times between these orders are shorter than the typical 15-45 second range
because I submitted several in close proximity.

```
2018/10/04 17:49:30  9999993   13.462s  28.6s    2018/10/04 17:52:22  2m51.598s
2018/10/04 17:49:44  9999994   13.525s  28.43s   2018/10/04 17:52:06  2m22.152s
2018/10/04 17:50:13  9999995   29.642s  28.444s  2018/10/04 17:52:07  1m53.775s
2018/10/04 17:50:48  9999996   34.984s  28.516s  2018/10/04 17:52:14  1m25.547s
2018/10/04 17:50:56  9999997   8.041s   28.293s  2018/10/04 17:51:53  56.586s
2018/10/04 17:50:59  9999998   3.059s   28.022s  2018/10/04 17:51:27  28.022s
2018/10/04 17:51:00  9999999   1.057s   27.735s  2018/10/04 17:51:00  0s
2018/10/04 17:51:01  10000000  1.045s   27.454s  2018/10/04 17:50:34  -27.454s
2018/10/04 17:51:02  10000001  1.01s    27.179s  2018/10/04 17:50:08  -54.357s
2018/10/04 17:51:03  10000002  409ms    26.903s  2018/10/04 17:49:42  -1m20.708s
2018/10/04 17:51:03  10000003  479ms    26.633s  2018/10/04 17:49:17  -1m46.532s
```
*Times are in MDT*
