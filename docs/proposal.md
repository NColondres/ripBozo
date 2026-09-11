# <img src="imgs/RIP_Bozo.webp" width="40"> RipBozo Proposal

**Put your money where your mouth is... or rather, your time.**

## 🎯 The Vision
Twitch chat is the wild west of arguments. **RipBozo** turns those arguments into high-stakes wagers. Users can challenge each other to a simple game where the stakes are their freedom in the channel. If you lose, you get banned for the duration of the wager.

---

## 🛠️ How it Works: The Cycle of Chaos

1. **The Challenge** ⚔️
   A user sends a command in chat: `!challenge @User 5m`. 
   *They are challenging @User to a game with a 5-minute ban on the line.*

2. **The Acceptance** ✅
   If the challenged user accepts, the RipBozo bot generates a unique, one-time link to a simple game.

3. **The Game** 🎮
   Both users jump into the RipBozo webapp. They play a quick, fair, and simple game (e.g., Rock-Paper-Scissors or a Coin Flip).

4. **The Rip** 💀
   The service determines the loser. The bot then:
   - **Bans the loser** for the wagered amount of time.
   - **Notifies the chat** with a victory message: 
     *"@Winner just RipBozo'd @Loser! See you in 5 minutes! <img src="imgs/RIP_Bozo.webp" width="20">"*

---

## 🛡️ Safeguards
To prevent total chat annihilation and keep the chaos manageable, RipBozo implements a few key rules:
- **One at a Time:** Only one challenge can be active in a channel at any given moment. No queueing, no multitasking—just one epic showdown.
- **The Cooldown:** A cooldown period is enforced between challenges. This ensures the channel doesn't lose half its population in the first five minutes.

---

## 🏗️ Technical Breakdown

The project is split into four main components:

### 🤖 Twitch Bot
- Parses channel messages in real-time.
- Handles the `!challenge` command and acceptance logic.
- Broadcasts the final result and the "Rip" notification to the chat.

### ⚙️ Challenge Service
- Manages the state of active challenges.
- Generates secure, temporary game links.
- Acts as the source of truth for who won and who lost.

### 🌐 RipBozo Webapp
- A lightweight frontend where the actual battle happens.
- Implements simple game logic (Rock-Paper-Scissors, Coin Flip).
- Reports the outcome back to the Challenge Service.

### 🔨 Ban Execution
- Interfaces with the Twitch API to apply the time-limited ban to the loser.
- Ensures the punishment is swift and accurate.

---

## 🚀 Future Roadmap
- [ ] **More Games:** Expand beyond RPS and Coin Flips to more advanced mini-games.
- [ ] **Leaderboards:** Track who the biggest "Bozo-Slayers" are in a channel.
- [ ] **Custom Wagers:** Allow streamers to set minimum/maximum wager limits.
- [ ] **Integration:** More flashy chat notifications and overlays.
