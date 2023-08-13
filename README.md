# gobuffer
A buffered file writer, Improving performance with dual memory caches.  
gobuffer是一个文件写入组件，通过双缓存的设计，在缓存满的时候切换活动缓存，减少锁的持有时间。最大化提高写入性能。
