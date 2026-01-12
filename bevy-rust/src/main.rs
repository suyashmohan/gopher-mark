use bevy::{
    asset::AssetMetaCheck,
    diagnostic::{DiagnosticsStore, FrameTimeDiagnosticsPlugin},
    prelude::*,
};
use rand::Rng;

const SCREEN_WIDTH: u32 = 1920;
const SCREEN_HEIGHT: u32 = 1080;
const GOPHER_COUNT: usize = 1000;
const FIXED_TIMESTEP_HZ: f64 = 60.0; // 60 ticks per second

fn main() {
    App::new()
        .add_plugins((
            DefaultPlugins
                .set(WindowPlugin {
                    primary_window: Some(Window {
                        title: "gopher mark".into(),
                        resolution: (SCREEN_WIDTH, SCREEN_HEIGHT).into(),
                        present_mode: bevy::window::PresentMode::AutoNoVsync,
                        ..default()
                    }),
                    ..default()
                })
                .set(AssetPlugin {
                    meta_check: AssetMetaCheck::Never,
                    ..default()
                }),
            FrameTimeDiagnosticsPlugin::default(),
        ))
        .insert_resource(Time::<Fixed>::from_hz(FIXED_TIMESTEP_HZ))
        .add_systems(Startup, setup)
        .add_systems(FixedUpdate, gopher_movement)
        .add_systems(Update, (spawn_gophers, update_fps_text))
        .run();
}

#[derive(Component)]
struct Gopher {
    pub velocity: Vec2,
}

#[derive(Component)]
struct FpsText;

#[derive(Resource)]
struct FpsUpdateTimer {
    timer: Timer,
}

#[derive(Resource)]
struct SpawnCounter {
    counter: usize,
}

#[derive(Resource)]
struct GopherAssets {
    handle: Handle<Image>,
}

fn setup(mut commands: Commands, asset_server: Res<AssetServer>) {
    let gopher_image = asset_server.load("gopher.png");

    commands.spawn(Camera2d);
    commands.spawn((
        Text::new("FPS: "),
        FpsText,
        Node {
            position_type: PositionType::Absolute,
            top: Val::Px(10.0),
            left: Val::Px(10.0),
            ..default()
        },
    ));

    commands.insert_resource(GopherAssets {
        handle: gopher_image.clone(),
    });
    commands.insert_resource(SpawnCounter { counter: 0 });
    commands.insert_resource(FpsUpdateTimer {
        timer: Timer::from_seconds(0.5, TimerMode::Repeating),
    });
}

fn spawn_gophers(
    mut commands: Commands,
    mut spawn_counter: ResMut<SpawnCounter>,
    assets: Res<GopherAssets>,
    inputs: Res<ButtonInput<KeyCode>>,
) {
    if inputs.just_released(KeyCode::Space) {
        let mut rng = rand::rng();
        let mut gophers_to_spawn = Vec::with_capacity(GOPHER_COUNT);

        for _ in 0..GOPHER_COUNT {
            let velocity = Vec2 {
                x: rng.random_range(-100.0..100.0),
                y: rng.random_range(-100.0..100.0),
            };
            gophers_to_spawn.push((
                Gopher { velocity },
                Sprite::from_image(assets.handle.clone()),
                Transform::from_xyz(0., 0., 0.),
            ));
        }

        commands.spawn_batch(gophers_to_spawn);
        spawn_counter.counter += GOPHER_COUNT;
    }
}

fn gopher_movement(time: Res<Time<Fixed>>, mut query: Query<(&mut Gopher, &mut Transform)>) {
    let half_width = SCREEN_WIDTH as f32 / 2.0;
    let half_height = SCREEN_HEIGHT as f32 / 2.0;

    let delta = time.delta_secs();

    for (mut gopher, mut transform) in query.iter_mut() {
        if transform.translation.x < -half_width || transform.translation.x > half_width {
            gopher.velocity.x *= -1.0;
        }
        if transform.translation.y < -half_height || transform.translation.y > half_height {
            gopher.velocity.y *= -1.0;
        }
        transform.translation.x += gopher.velocity.x * delta;
        transform.translation.y += gopher.velocity.y * delta;
    }
}

fn update_fps_text(
    diagnostics: Res<DiagnosticsStore>,
    spawn_counter: Res<SpawnCounter>,
    time: Res<Time>,
    mut fps_update_timer: ResMut<FpsUpdateTimer>,
    mut query: Query<&mut Text, With<FpsText>>,
) {
    fps_update_timer.timer.tick(time.delta());

    if fps_update_timer.timer.just_finished() {
        if let Some(fps) = diagnostics.get(&FrameTimeDiagnosticsPlugin::FPS) {
            if let Some(value) = fps.smoothed() {
                for mut text in &mut query {
                    text.0 = format!("FPS: {:.1} - Gophers: {}", value, spawn_counter.counter);
                }
            }
        }
    }
}
