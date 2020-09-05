import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';

import { RecipeModule } from './recipe/recipe-module'

@Module({
    imports: [RecipeModule,
        TypeOrmModule.forRoot({
            type: 'postgres',
            host: 'localhost',
            port: 5432,
            username: 'peppermint-user',
            password: '1234',
            database: 'peppermint-db',
            entities: ["dist/**/*-model*"],
            synchronize: true,
        }),
    ],
    controllers: [],
    providers: [],
})
export class AppModule {}
