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
   If the challenged user accepts, the bot updates the database and generates unique, one-time links for both participants. These links (e.g., `https://www.ripbozo.com/challenge?challenge_id=foo&player_id=bar`) are sent directly to the users in the Twitch chat.

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
- **URL Handoff:** Upon acceptance, generates unique, secure links for each player (e.g., `https://www.ripbozo.com/challenge?challenge_id=foo&player_id=bar`) and sends them to the users in chat.
- Broadcasts the final result and the "Rip" notification to the chat after receiving the outcome from the webapp.

### 🌐 Next.js Application (Frontend & Backend)
- **Identity Verification:** Parses the `challenge_id` and `player_id` from the URL, verifying the participant against the PostgreSQL database to ensure they are authorized for that specific challenge.
- **Game Management:** Displays the correct game, identifies the user, and presents the stakes (the ban duration) retrieved from the DB.
- **Game Logic:** Executes the mini-game (RPS, Coin Flip) and determines the winner.
- **Outcome Reporting:** Once a winner is declared, sends a POST request to the Twitch Bot with the final result to trigger the punishment.

### 🔨 Ban Execution
- Triggered by the Twitch Bot upon receiving the outcome from the Next.js app.
- Interfaces with the Twitch API to apply the time-limited ban to the loser.
- Ensures the punishment is swift and accurate.

---

## 🎨 Frontend Experience
To match the high-stakes, chaotic energy of Twitch, the webapp will focus on a "Digital Arcade" aesthetic—dark mode by default with high-contrast neon accents.

### 🕹️ Game Rendering & UX
- **Dynamic Game States:** The UI will transition smoothly between states:
    - `Waiting`: A suspenseful screen showing the opponent's name and a "Waiting for opponent to join..." spinner.
    - `Active`: The main game interface (e.g., three large, glowing buttons for RPS) with a real-time countdown timer.
    - `Resolution`: A dramatic reveal of both players' choices.
- **Stakes Visibility:** The ban duration (e.g., "5 MINUTES ON THE LINE") will be persistently displayed in a bold, warning-style banner to maintain tension.
- **Responsive Design:** Optimized for mobile browsers, as most Twitch users will click the link from their phones.

### ✨ Animations & "Juice"
To make the experience feel polished and exciting, we'll implement:
- **The Build-up:** Use a "3... 2... 1..." countdown animation before the final result is revealed to create a peak moment of suspense.
- **Visual Feedback:** 
    - **Victory:** Confetti explosions and glowing gold borders for the winner.
    - **The Rip:** A "glitch" or "screen tear" animation for the loser, visually simulating their "removal" from the chat.
- **Smooth Transitions:** Using libraries like **Framer Motion** for fluid entry/exit animations and button hover effects.

---

## 🚀 Future Roadmap
- [ ] **More Games:** Expand beyond RPS and Coin Flips to more advanced mini-games.
- [ ] **Leaderboards:** Track who the biggest "Bozo-Slayers" are in a channel.
- [ ] **Custom Wagers:** Allow streamers to set minimum/maximum wager limits.
- [ ] **Integration:** More flashy chat notifications and overlays.
