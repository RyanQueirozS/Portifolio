# Handeling FPS on SDL

I've had [this project](https://github.com/RyanQueirozS/SDL_cpp) for some time
now, it's supposed to be a simple game engine and this topic has by far been
the hardest/least unintuitive

## The problem

Games deal with FPS, and most have a way to cap it. FPS, as it entails, is how
many frames happen in a second. If, for example, FPS is capped at 60, it means
that in 1 second, there should be a maximum of 60 frames, simplifying **60fps = 1
second/60 frames**. 

In our case, we will use milliseconds (1s = 1000ms) for code simplicity.

## Developing the problem

```cpp
void Engine::HandleFps() {
    float timeStepInMs = 1000.0f / fpsCap;
}
```
In this case, **fpsCap** is defined in the **Engine** class and its value is 60.

Secondly, we need to have access to the **first tick of the loop**, we can
simply add a float as a param and **compare it to the last tick** of the loop

```cpp
void Engine::HandleFps(float startTick) {
    float timeStepInMs = 1000.0f / fpsCap;
    float lastLoopTick = SDL_GetTicks();
    float timeElapsed = lastLoopTick - startTick; // I won't use this later, it's just for demonstration
}
```

Lastly, we need to delay the game 'til it reaches the desired frames in a
second. And this is where it gets unintuitive.

My confusion was: 
- Wouldn't that worsen the game's performance? 
- Why would I implement that instead of just always trying to reach the maximum FPS?

>"Wouldn't that worsen the game's performance?"

Like always, it depends. Because this project is super small and barely uses
any hardware, maybe? The tests I've performed didn't show sufficient difference.

>"Why would I implement that instead of just always trying to reach the maximum FPS?"

This is the opposite of what was before. In regards to computer heavy games,
having the ability to cap resources is great. Having a constant FPS is better
than having a lot of FPS sometimes.

Said that we still need to delay the game, and most importantly, check if the fps is slower than expected

```cpp
    if (timeStepInMs > timeElapsed) {
        sdl_delay(timeStepInMs - timeElapsed);
        return;
    }
```

## Final code

```cpp
void Engine::HandleFPS(float startTick) {
    float timeStepInMS = 1000.0f / fpsCap;
    float lastLoopTick = SDL_GetTicks();
    if (timeStepInMs > (lastLoopTick - startTick)) {
        sdl_delay(timeStepInMs - lastLoopTick - startTick);
        return;
    }
}
```
