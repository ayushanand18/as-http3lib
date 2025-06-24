namespace ashttp3lib :: utility {
    //! \brief MessageQueue class. A Concurrent Thread-safe implementation
    //! message queue to store pending packets.
    template <typename T>
    class MessageQueue {
        //! \brief Job Pool. A pool of pending jobs that must be executed. 
        std::queue<T> jobPool;

        //! \brief Push a job. Push a job into the queue to be 
        //! executed asynchronously.
        //! \param newJob The new job submitted into the queue.
        async bool push(T newJob) noexcept {
            // mutex.lock
            // insert into the queue
            // mutex.unlock
        }

        //! \brief Pop a job. Pop a job from the queue to be 
        //! executed asynchronously.
        async T pop() {
            // mutex.lock
            // insert into the queue
            // mutex.unlock
        }

    }; // class MessageQueue
} // namespace ashttp3lib :: utility