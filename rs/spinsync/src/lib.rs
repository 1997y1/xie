// const MAX_BACKOFF: u32 = 8;

pub struct Spinlock {
    state: std::sync::atomic::AtomicU32,
}

impl Spinlock {
    pub const fn new() -> Self {
        Spinlock {
            state: std::sync::atomic::AtomicU32::new(0),
        }
    }

    pub fn lock(&self) {
        // let mut backoff = 0;

        while self
            .state
            .compare_exchange(
                0,
                1,
                std::sync::atomic::Ordering::Acquire,
                std::sync::atomic::Ordering::Relaxed,
            )
            .is_err()
        {
            // Self::auto_schedule(&mut backoff);// Threads need to be enabled.
        }
    }

    pub fn unlock(&self) {
        self.state.store(0, std::sync::atomic::Ordering::Release);
    }

    // fn auto_schedule(backoff: &mut u32) {
    //     *backoff += 1;
    //     if *backoff > MAX_BACKOFF {
    //         *backoff = 0;
    //         std::thread::yield_now();
    //     }
    // }
}
