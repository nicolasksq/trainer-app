# Usage Guide

How to get the most out of your AI endurance coach.

## Getting Started

1. Open Claude Code in the `trainer-app` directory
2. Type `@trainer` followed by your request
3. The agent starts both MCP servers automatically -- no manual setup needed

For your first conversation, let the coach learn about you:

```
@trainer I'm a 35-year-old runner, 3 years of experience. I run 4 days/week
and I'm training for a sub-1:45 half marathon in October.
```

The agent will save your profile and use it to personalize all future recommendations.

## Example Prompts by Category

### Training Analysis

```
@trainer Analyze my last 2 weeks of training
@trainer Review my long run from Sunday -- was I in the right zones?
@trainer Compare my running volume this month vs last month
@trainer What's my current estimated 10K time?
@trainer Show me my pace and heart rate trends over the last month
```

### Training Plans

```
@trainer Create a 12-week plan to run a sub-40 10K
@trainer I have a half marathon in 8 weeks, build me a plan
@trainer Adjust my plan -- I can only train 3 days this week
@trainer I swim Mon/Wed/Fri, schedule running around that
@trainer Move me to a recovery week, I'm feeling fatigued
```

### Workout Scheduling

```
@trainer Schedule this week's workouts to my Garmin
@trainer Create an interval session: 6x1km at 4:00/km with 2min recovery
@trainer Program a tempo run for Thursday on my watch
@trainer Show me the workouts currently on my Garmin calendar
@trainer Delete last week's unfinished workouts from Garmin
```

### Recovery and Health

```
@trainer How's my recovery today?
@trainer Check my HRV and sleep trends this week
@trainer Am I overtraining? Look at my training load
@trainer What does my Garmin training readiness say?
@trainer Compare my resting heart rate over the last 2 weeks
```

### Goal Setting

```
@trainer I want to run a sub-1:45 half marathon by October
@trainer Help me improve my cycling FTP
@trainer I'm training for my first triathlon in 3 months
@trainer Set a goal to increase my weekly mileage to 60km
```

## How Memory Works

The `@trainer` agent remembers your information across sessions using Claude Code's memory system. It persists:

- **Athlete profile** -- your experience level, sport history, injuries, and physical data
- **Current training plan** -- the active plan, current phase (base, build, peak, taper), and weekly structure
- **Goals and events** -- target races, goal times, and deadlines
- **Preferences** -- available training days, schedule constraints, equipment, and communication style

The agent reads this context at the start of each conversation to maintain continuity. You do not need to repeat your background every time.

To explicitly save or update your profile:

```
@trainer Save my profile: I now have access to a pool and want to add swimming
@trainer Update my goal -- I changed my target race to a marathon in March
```

## Multi-Sport Support

The agent handles multiple sports:

- **Running** -- road, trail, track workouts with pace and HR targets
- **Cycling** -- endurance rides, intervals, FTP work with power targets
- **Swimming** -- pool sessions with distance and pace guidance
- **Strength training** -- gym sessions with exercises, sets, reps, and RPE

Data is read from both Strava (activities, GPS streams, power data) and Garmin (health metrics, recovery scores, body composition). Structured workouts can be pushed to your Garmin device for any supported sport.

## Tips for Best Results

- **Be specific about constraints** -- share your available days, time per session, equipment access, and other commitments
- **Share your race calendar** -- the agent builds plans around your target events and goal times
- **Ask for post-workout feedback** -- after key sessions (long runs, intervals, tempo), ask the agent to review your execution against the prescription
- **Check in weekly** -- ask for plan adjustments based on how the week went, recovery status, and upcoming schedule
- **Report how you feel** -- fatigue, soreness, motivation, and life stress all factor into smart training adjustments
- **All paces are in min/km** -- the agent displays pace as min:sec per kilometer (e.g., 4:30/km), never in km/h or m/s
