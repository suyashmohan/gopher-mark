/*
Raylib example file.
This is an example main file for a simple raylib project.
Use this as a starting point or replace it with your code.

by Jeffery Myers is marked with CC0 1.0. To view a copy of this license, visit https://creativecommons.org/publicdomain/zero/1.0/

*/

#include <stdlib.h>
#include <time.h>
#include "raylib.h"

#include "resource_dir.h" // utility header for SearchAndSetResourceDir

typedef struct
{
	int PosX;
	int PosY;
	int VelX;
	int VelY;
} gopher_t;

gopher_t gopher_new(int maxX, int maxY)
{
	gopher_t g;
	g.PosX = rand() % maxX;
	g.PosY = rand() % maxY;
	g.VelX = 1 + (rand() % 10);
	g.VelY = 1 + (rand() % 10);
	return g;
}

void gopher_move(gopher_t *g, int w, int h, int maxX, int maxY)
{
	if (g->PosX + g->VelX > maxX - w || g->PosX < 0)
	{
		g->VelX *= -1;
	}

	if (g->PosY + g->VelY > maxY - h || g->PosY < 0)
	{
		g->VelY *= -1;
	}

	g->PosX += g->VelX;
	g->PosY += g->VelY;
}

int main()
{
	const int screenWidth = 1920;
	const int screenHeight = 1080;

	srand(time(NULL));

	InitWindow(screenWidth, screenHeight, "gopher mark");
	SetTargetFPS(60);

	// Utility function from resource_dir.h to find the resources folder and set it as the current working directory so we can load from it
	SearchAndSetResourceDir("resources");

	Texture gopher_texture = LoadTexture("gopher.png");

	int gopherAddCount = 5000;
	int gopherCapacity = gopherAddCount; // Initial Capacity
	int gopherSize = 0;
	gopher_t *gophers = (gopher_t *)malloc(gopherCapacity * sizeof(gopher_t));
	for (int i = 0; i < gopherAddCount; i++)
	{
		gophers[gopherSize++] = gopher_new(screenWidth, screenHeight);
	}

	while (!WindowShouldClose()) // run the loop until the user presses ESCAPE or presses the Close button on the window
	{
		if (IsKeyReleased(KEY_SPACE))
		{
			gopherCapacity += gopherAddCount;
			gophers = (gopher_t *)realloc(gophers, gopherCapacity * sizeof(gopher_t));
			for (int i = 0; i < gopherAddCount; i++)
			{
				gophers[gopherSize++] = gopher_new(screenWidth, screenHeight);
			}
		}

		for (int i = 0; i < gopherSize; i++)
		{
			gopher_move(&gophers[i], gopher_texture.width, gopher_texture.height, screenWidth, screenHeight);
		}

		BeginDrawing();
		ClearBackground(WHITE);

		for (int i = 0; i < gopherSize; i++)
		{
			DrawTexture(gopher_texture, gophers[i].PosX, gophers[i].PosY, WHITE);
		}

		DrawText(TextFormat("Gophers: %d", gopherSize), 8, 8, 20, BLACK);
		DrawFPS(8, 30);

		EndDrawing();
	}

	// cleanup
	UnloadTexture(gopher_texture);
	free(gophers);
	CloseWindow();
	return 0;
}
